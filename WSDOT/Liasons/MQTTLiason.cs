using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using TwoMQTT;
using TwoMQTT.Interfaces;
using TwoMQTT.Liasons;
using TwoMQTT.Models;
using TwoMQTT.Utils;
using WSDOT.Models.Options;
using WSDOT.Models.Shared;

namespace WSDOT.Liasons;

/// <summary>
/// An class representing a managed way to interact with MQTT.
/// </summary>
public class MQTTLiason : MQTTLiasonBase<Resource, object, SlugMapping, SharedOpts>, IMQTTLiason<Resource, object>
{
    /// <summary>
    /// 
    /// </summary>
    /// <param name="logger"></param>
    /// <param name="generator"></param>
    /// <param name="sharedOpts"></param>
    public MQTTLiason(ILogger<MQTTLiason> logger, IMQTTGenerator generator, IOptions<SharedOpts> sharedOpts) :
        base(logger, generator, sharedOpts)
    {
    }

    /// <inheritdoc />
    public IEnumerable<(string topic, string payload)> MapData(Resource input)
    {
        var results = new List<(string, string)>();
        var slug = this.Questions
            .Where(x => x.TravelTimeID == input.TravelTimeID)
            .Select(x => x.Slug)
            .FirstOrDefault() ?? string.Empty;

        this.Logger.LogDebug("Found slug {slug} for incoming data for {ttid}", slug, input.TravelTimeID);
        if (string.IsNullOrEmpty(slug))
        {
            this.Logger.LogDebug("Unable to find slug for {ttid}", input.TravelTimeID);
            return results;
        }

        this.Logger.LogDebug("Found slug {slug} for incoming data for {ttid}", slug, input.TravelTimeID);
        results.AddRange(new[]
            {
                    (this.Generator.StateTopic(slug), input.CurrentTime.ToString()),
                }
        );

        return results;
    }

    /// <inheritdoc />
    public IEnumerable<(string slug, string sensor, string type, MQTTDiscovery discovery)> Discoveries()
    {
        var discoveries = new List<(string, string, string, MQTTDiscovery)>();
        var assembly = Assembly.GetAssembly(typeof(Program))?.GetName() ?? new AssemblyName();
        var mapping = new[]
        {
                new { Sensor = string.Empty, Type = Const.SENSOR },
            };

        foreach (var input in this.Questions)
        {
            foreach (var map in mapping)
            {
                this.Logger.LogDebug("Generating discovery for {ttid} - {sensor}", input.TravelTimeID, map.Sensor);
                var discovery = this.Generator.BuildDiscovery(input.Slug, map.Sensor, assembly, false);
                discovery = discovery with
                {
                    Icon = "mdi:car",
                    UnitOfMeasurement = "min",
                };

                discoveries.Add((input.Slug, map.Sensor, map.Type, discovery));
            }
        }

        return discoveries;
    }
}
