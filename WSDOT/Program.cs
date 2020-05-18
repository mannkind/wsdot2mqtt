using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using TwoMQTT.Core;
using TwoMQTT.Core.DataAccess;
using TwoMQTT.Core.Extensions;
using WSDOT.DataAccess;
using WSDOT.Managers;
using WSDOT.Models.Shared;


namespace WSDOT
{
    class Program : ConsoleProgram
    {
        static async Task Main(string[] args)
        {
            var p = new Program();
            p.MapOldEnvVariables();
            await p.ExecuteAsync(args);
        }

        protected override IServiceCollection ConfigureServices(HostBuilderContext hostContext, IServiceCollection services)
        {
            var sharedSect = hostContext.Configuration.GetSection(Models.Shared.Opts.Section);
            var sourceSect = hostContext.Configuration.GetSection(Models.SourceManager.Opts.Section);
            var sinkSect = hostContext.Configuration.GetSection(Models.SinkManager.Opts.Section);

            services.AddHttpClient<IHTTPSourceDAO<SlugMapping, Command, Models.SourceManager.FetchResponse, object>>();

            return services
                .Configure<Models.Shared.Opts>(sharedSect)
                .Configure<Models.SourceManager.Opts>(sourceSect)
                .Configure<Models.SinkManager.Opts>(sinkSect)
                .AddTransient<IHTTPSourceDAO<SlugMapping, Command, Models.SourceManager.FetchResponse, object>>(x =>
                {
                    var opts = x.GetService<IOptions<Models.SourceManager.Opts>>();
                    return new SourceDAO(x.GetService<ILogger<SourceDAO>>(), x.GetService<IHttpClientFactory>(), opts.Value.ApiKey);
                })
                .ConfigureBidirectionalSourceSink<Resource, Command, SourceManager, SinkManager>();
        }

        [Obsolete("Remove in the near future.")]
        private void MapOldEnvVariables()
        {
            var found = false;
            var foundOld = new List<string>();
            var mappings = new[]
            {
                new { Src = "WSDOT_SECRET", Dst = "WSDOT__APIKEY", CanMap = true, Strip = "", Sep = "" },
                new { Src = "WSDOT_TRAVELTIMEMAPPING", Dst = "WSDOT__RESOURCES", CanMap = true, Strip = "",  Sep = ":" },
                new { Src = "WSDOT_LOOKUPINTERVAL", Dst = "WSDOT__POLLINGINTERVAL", CanMap = false, Strip = "", Sep = "" },
                new { Src = "MQTT_TOPICPREFIX", Dst = "WSDOT__MQTT__TOPICPREFIX", CanMap = true, Strip = "", Sep = "" },
                new { Src = "MQTT_DISCOVERY", Dst = "WSDOT__MQTT__DISCOVERYENABLED", CanMap = true, Strip = "", Sep = "" },
                new { Src = "MQTT_DISCOVERYPREFIX", Dst = "WSDOT__MQTT__DISCOVERYPREFIX", CanMap = true, Strip = "", Sep = "" },
                new { Src = "MQTT_DISCOVERYNAME", Dst = "WSDOT__MQTT__DISCOVERYNAME", CanMap = true, Strip = "", Sep = "" },
                new { Src = "MQTT_BROKER", Dst = "WSDOT__MQTT__BROKER", CanMap = true, Strip = "tcp://", Sep = "" },
                new { Src = "MQTT_USERNAME", Dst = "WSDOT__MQTT__USERNAME", CanMap = true, Strip = "", Sep = "" },
                new { Src = "MQTT_PASSWORD", Dst = "WSDOT__MQTT__PASSWORD", CanMap = true, Strip = "", Sep = "" },
            };

            foreach (var mapping in mappings)
            {
                var old = Environment.GetEnvironmentVariable(mapping.Src);
                if (string.IsNullOrEmpty(old))
                {
                    continue;
                }

                found = true;
                foundOld.Add($"{mapping.Src} => {mapping.Dst}");

                if (!mapping.CanMap)
                {
                    continue;
                }

                // Strip junk where possible
                if (!string.IsNullOrEmpty(mapping.Strip))
                {
                    old = old.Replace(mapping.Strip, string.Empty);
                }

                // Simple
                if (string.IsNullOrEmpty(mapping.Sep))
                {
                    Environment.SetEnvironmentVariable(mapping.Dst, old);
                }
                // Complex
                else
                {
                    var resourceSlugs = old.Split(",");
                    var i = 0;
                    foreach (var resourceSlug in resourceSlugs)
                    {
                        var parts = resourceSlug.Split(mapping.Sep);
                        var id = parts.Length >= 1 ? parts[0] : string.Empty;
                        var slug = parts.Length >= 2 ? parts[1] : string.Empty;
                        var idEnv = $"{mapping.Dst}__{i}__TravelTimeID";
                        var slugEnv = $"{mapping.Dst}__{i}__Slug";
                        Environment.SetEnvironmentVariable(idEnv, id);
                        Environment.SetEnvironmentVariable(slugEnv, slug);
                    }
                }

            }


            if (found)
            {
                var loggerFactory = LoggerFactory.Create(builder => { builder.AddConsole(); });
                var logger = loggerFactory.CreateLogger<Program>();
                logger.LogWarning("Found old environment variables.");
                logger.LogWarning($"Please migrate to the new environment variables: {(string.Join(", ", foundOld))}");
                Thread.Sleep(5000);
            }
        }
    }
}
