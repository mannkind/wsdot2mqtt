using System.Collections.Generic;
using TwoMQTT.Core.Interfaces;
using WSDOT.Models.Shared;

namespace WSDOT.Models.Options
{
    /// <summary>
    /// The shared options across the application
    /// </summary>
    public record SharedOpts : ISharedOpts<SlugMapping>
    {
        public const string Section = "WSDOT";

        /// <summary>
        /// 
        /// </summary>
        /// <typeparam name="SlugMapping"></typeparam>
        /// <returns></returns>
        public List<SlugMapping> Resources { get; init; } = new List<SlugMapping>();
    }
}
