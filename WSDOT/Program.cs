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
    class Program
    {
        static async Task Main(string[] args) => await ConsoleProgram<Resource, object, SourceLiason, MQTTLiason>.
            ExecuteAsync(args,
                envs: new Dictionary<string, string>()
                {
                    {
                        $"{Models.Options.MQTTOpts.Section}:{nameof(Models.Options.MQTTOpts.TopicPrefix)}",
                        Models.Options.MQTTOpts.TopicPrefixDefault
                    },
                    {
                        $"{Models.Options.MQTTOpts.Section}:{nameof(Models.Options.MQTTOpts.DiscoveryName)}",
                        Models.Options.MQTTOpts.DiscoveryNameDefault
                    },
                },
                configureServices: (HostBuilderContext context, IServiceCollection services) =>
                {
                    services
                        .AddOptions<Models.Options.SharedOpts>(Models.Options.SharedOpts.Section, context.Configuration)
                        .AddOptions<Models.Options.SourceOpts>(Models.Options.SourceOpts.Section, context.Configuration)
                        .AddOptions<TwoMQTT.Models.MQTTManagerOptions>(Models.Options.MQTTOpts.Section, context.Configuration)
                        .AddHttpClient()
                        .AddSingleton<IThrottleManager, ThrottleManager>(x =>
                        {
                            var opts = x.GetRequiredService<IOptions<Models.Options.SourceOpts>>();
                            return new ThrottleManager(opts.Value.PollingInterval);
                        })
                        .AddSingleton<ISourceDAO, SourceDAO>(x =>
                        {
                            var logger = x.GetRequiredService<ILogger<SourceDAO>>();
                            var httpClientFactory = x.GetRequiredService<IHttpClientFactory>();
                            var opts = x.GetRequiredService<IOptions<Models.Options.SourceOpts>>();
                            return new SourceDAO(logger, httpClientFactory, opts.Value.ApiKey);
                        });
                });
    }
}