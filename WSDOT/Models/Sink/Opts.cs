using TwoMQTT.Core.Models;

namespace WSDOT.Models.SinkManager
{
    /// <summary>
    /// The sink options
    /// </summary>
    public class Opts : MQTTManagerOptions
    {
        public const string Section = "WSDOT:MQTT";

        /// <summary>
        /// 
        /// </summary>
        public Opts()
        {
            this.TopicPrefix = "home/wsdot";
            this.DiscoveryName = "wsdot";
        }
    }
}
