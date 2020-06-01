using System;

namespace WSDOT.Models.SourceManager
{
    /// <summary>
    /// The source options
    /// </summary>
    public class Opts
    {
        public const string Section = "WSDOT";

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public string ApiKey { get; set; } = string.Empty;

        /// <summary>
        /// 
        /// </summary>
        /// <returns></returns>
        public TimeSpan PollingInterval { get; set; } = new TimeSpan(0, 3, 31);
    }
}
