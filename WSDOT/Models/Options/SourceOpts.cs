using System;

namespace WSDOT.Models.Options
{
    /// <summary>
    /// The source options
    /// </summary>
    public record SourceOpts
    {
        public const string Section = "WSDOT";

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public string ApiKey { get; init; } = string.Empty;

        /// <summary>
        /// 
        /// </summary>
        /// <returns></returns>
        public TimeSpan PollingInterval { get; init; } = new TimeSpan(0, 3, 31);
    }
}
