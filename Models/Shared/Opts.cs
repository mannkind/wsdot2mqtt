using System.Collections.Generic;

namespace WSDOT.Models.Shared
{
    /// <summary>
    /// The shared options across the application
    /// </summary>
    public class Opts
    {
        public const string Section = "WSDOT:Shared";

        public List<SlugMapping> Resources { get; set; } = new List<SlugMapping>();
    }
}
