using System.Collections.Generic;
using System.Threading.Channels;
using System.Threading.Tasks;
using System.Threading;
using System.Linq;
using TwoMQTT.Core.Managers;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using Newtonsoft.Json;
using System.Net.Http;
using WSDOT.Models.Shared;

namespace WSDOT.Managers
{
    public class Source : HTTPManager<SlugMapping, Resource, Command>
    {
        public Source(ILogger<Source> logger, IOptions<Opts> sharedOpts, IOptions<Models.SourceManager.Opts> opts, ChannelWriter<Resource> outgoing, ChannelReader<Command> incoming, IHttpClientFactory httpClientFactory) :
            base(logger, outgoing, incoming, httpClientFactory.CreateClient())
        {
            this.opts = opts.Value;
            this.sharedOpts = sharedOpts.Value;
        }
        protected readonly Models.SourceManager.Opts opts;
        protected readonly Opts sharedOpts;

        /// <inheritdoc />
        protected override void LogSettings()
        {
            var resources = string.Join(",",
                this.sharedOpts.Resources.Select(x => $"{x.TravelTimeID}:{x.Slug}")
            );

            this.logger.LogInformation(
                $"ApiKey:                {this.opts.ApiKey}\n" +
                $"PollingInterval:       {this.opts.PollingInterval}\n" +
                $"Resources:             {resources}\n" +
                $""
            );
        }

        /// <inheritdoc />
        protected override async Task PollAsync(CancellationToken cancellationToken = default)
        {
            this.logger.LogInformation("Polling");

            var tasks = new List<Task<Models.SourceManager.Response>>();
            foreach (var key in this.sharedOpts.Resources)
            {
                this.logger.LogInformation($"Looking up {key}");
                tasks.Add(this.FetchOneAsync(key, cancellationToken));
            }

            var results = await Task.WhenAll(tasks);
            foreach (var result in results.Where(x => x.Ok))
            {
                this.logger.LogInformation($"Found {result}");
                await this.outgoing.WriteAsync(Resource.From(result), cancellationToken);
            }
        }

        /// <inheritdoc />
        protected override Task DelayAsync(CancellationToken cancellationToken = default) => 
            Task.Delay(this.opts.PollingInterval, cancellationToken);

        /// <summary>
        /// Fetch one record from the source
        /// </summary>
        private async Task<Models.SourceManager.Response> FetchOneAsync(SlugMapping key, CancellationToken cancellationToken = default)
        {
            var baseUrl = "https://www.wsdot.wa.gov/Traffic/api/TravelTimes/TravelTimesREST.svc/GetTravelTimeAsJson";
            var query = $"AccessCode={this.opts.ApiKey}&TravelTimeID={key.TravelTimeID}";
            var resp = await this.client.GetAsync($"{baseUrl}?{query}", cancellationToken);
            if (!resp.IsSuccessStatusCode)
            {
                return new Models.SourceManager.Response();
            }

            var content = await resp.Content.ReadAsStringAsync();
            var obj = JsonConvert.DeserializeObject<Models.SourceManager.Response>(content);
            obj.Ok = true;

            return obj;
        }
    }
}
