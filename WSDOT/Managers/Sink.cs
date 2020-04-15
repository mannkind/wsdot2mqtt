using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using System.Threading;
using System.Threading.Channels;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using MQTTnet;
using TwoMQTT.Core.Managers;
using WSDOT.Models.Shared;

namespace WSDOT.Managers
{
    public class Sink : MQTTManager<Resource, Command>
    {
        public Sink(ILogger<Sink> logger, IOptions<Opts> sharedOpts, IOptions<Models.SinkManager.Opts> opts, ChannelReader<Resource> inputChannel, ChannelWriter<Command> outputChannel) :
            base(logger, opts, inputChannel, outputChannel)
        {
            this.sharedOpts = sharedOpts.Value;
        }
        protected readonly Opts sharedOpts;

        /// <inheritdoc />
        protected override async Task HandleIncomingAsync(Resource input, CancellationToken cancellationToken = default)
        {
            var slug = this.sharedOpts.Resources
                .Where(x => x.TravelTimeID == input.TravelTimeID)
                .Select(x => x.Slug)
                .FirstOrDefault() ?? string.Empty;

            if (string.IsNullOrEmpty(slug))
            {
                return;
            }

            var topic = this.StateTopic(slug);
            var payload = input.CurrentTime.ToString();

            if (this.knownMessages.ContainsKey(topic) && this.knownMessages[topic] == payload)
            {
                this.logger.LogDebug($"Duplicate '{payload}' found on '{topic}'");
                return;
            }

            this.logger.LogInformation($"Publishing '{payload}' on '{topic}'");
            await this.client.PublishAsync(
                new MqttApplicationMessageBuilder()
                    .WithTopic(topic)
                    .WithPayload(payload)
                    .WithExactlyOnceQoS()
                    .WithRetainFlag()
                    .Build(),
                cancellationToken
            );

            this.knownMessages[topic] = payload;
        }

        /// <inheritdoc />
        protected override async Task HandleDiscoveryAsync(CancellationToken cancellationToken = default)
        {
            if (!this.opts.DiscoveryEnabled)
            {
                return;
            }

            var tasks = new List<Task>();
            var assembly = Assembly.GetAssembly(typeof(Program))?.GetName() ?? new AssemblyName();
            var mapping = new [] 
            {
                new { Sensor = string.Empty, Type = "sensor" },
            };

            foreach (var input in this.sharedOpts.Resources)
            {
                foreach (var map in mapping) 
                {
                    var discovery = this.BuildDiscovery(input.Slug, map.Sensor, assembly, false);
                    discovery.Icon = "mdi:car";
                    discovery.UnitOfMeasurement = "min";
                    tasks.Add(this.PublishDiscoveryAsync(input.Slug, map.Sensor, map.Type, discovery, cancellationToken));
                }
            }

            await Task.WhenAll(tasks);
        }
    }
}