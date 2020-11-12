using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Threading.Tasks;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using TwoMQTT;
using TwoMQTT.Extensions;
using TwoMQTT.Interfaces;
using TwoMQTT.Managers;
using WSDOT.DataAccess;
using WSDOT.Liasons;
using WSDOT.Models.Shared;

namespace WSDOT
{
    class Program : ConsoleProgram<Resource, object, SourceLiason, MQTTLiason>
    {
        static async Task Main(string[] args)
        {
            var p = new Program();
            await p.ExecuteAsync(args);
        }

        /// <inheritdoc />
        protected override IDictionary<string, string> EnvironmentDefaults()
        {
            var sep = "__";
            var section = Models.Options.MQTTOpts.Section.Replace(":", sep);
            var sectsep = $"{section}{sep}";

            return new Dictionary<string, string>
            {
                { $"{sectsep}{nameof(Models.Options.MQTTOpts.TopicPrefix)}", Models.Options.MQTTOpts.TopicPrefixDefault },
                { $"{sectsep}{nameof(Models.Options.MQTTOpts.DiscoveryName)}", Models.Options.MQTTOpts.DiscoveryNameDefault },
            };
        }

        /// <inheritdoc />
        protected override IServiceCollection ConfigureServices(HostBuilderContext hostContext, IServiceCollection services)
        {
            services.AddHttpClient<ISourceDAO>();

            return services
                .ConfigureOpts<Models.Options.SharedOpts>(hostContext, Models.Options.SharedOpts.Section)
                .ConfigureOpts<Models.Options.SourceOpts>(hostContext, Models.Options.SourceOpts.Section)
                .ConfigureOpts<TwoMQTT.Models.MQTTManagerOptions>(hostContext, Models.Options.MQTTOpts.Section)
                .AddSingleton<IThrottleManager, ThrottleManager>(x =>
                {
                    var opts = x.GetService<IOptions<Models.Options.SourceOpts>>();
                    if (opts == null)
                    {
                        throw new ArgumentException($"{nameof(opts.Value.PollingInterval)} is required for {nameof(ThrottleManager)}.");
                    }

                    return new ThrottleManager(opts.Value.PollingInterval);
                })
                .AddSingleton<ISourceDAO, SourceDAO>(x =>
                {
                    var logger = x.GetService<ILogger<SourceDAO>>();
                    var httpClientFactory = x.GetService<IHttpClientFactory>();
                    var opts = x.GetService<IOptions<Models.Options.SourceOpts>>();

                    if (logger == null)
                    {
                        throw new ArgumentException($"{nameof(logger)} is required for {nameof(SourceDAO)}.");
                    }
                    if (httpClientFactory == null)
                    {
                        throw new ArgumentException($"{nameof(httpClientFactory)} is required for {nameof(SourceDAO)}.");
                    }
                    if (opts == null)
                    {
                        throw new ArgumentException($"{nameof(opts.Value.ApiKey)} are required for {nameof(SourceDAO)}.");
                    }

                    return new SourceDAO(logger, httpClientFactory, opts.Value.ApiKey);
                });
        }
    }
}