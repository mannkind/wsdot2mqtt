using TwoMQTT.Core.Models;

namespace WSDOT.Models.Options
{
    /// <summary>
    /// The sink options
    /// </summary>
    public class MQTTOpts : MQTTManagerOptions
    {
        public const string Section = "WSDOT:MQTT";
        public const string TopicPrefixDefault = "home/wsdot";
        public const string DiscoveryNameDefault = "wsdot";
    }
}
