using System;

namespace WSDOT.Models.SourceManager
{
    /// <summary>
    /// The source options
    /// </summary>
    public class Opts
    {
        public const string Section = "WSDOT:Source";

        public string ApiKey { get; set; } = string.Empty;
        public TimeSpan PollingInterval { get; set; } = new TimeSpan(0, 3, 31);

        public override string ToString() => $"ApiKey: {this.ApiKey}, Polling Interval: {this.PollingInterval}";
    }
}
