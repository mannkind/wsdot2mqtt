namespace WSDOT.Models.Shared
{
    /// <summary>
    /// The shared key info => slug mapping across the application
    /// </summary>
    public record SlugMapping
    {
        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public long TravelTimeID { get; init; } = 0;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public string Slug { get; init; } = string.Empty;

        /// <inheritdoc />
        public override string ToString() => $"Time Travel ID: {this.TravelTimeID}, Slug: {this.Slug}";
    }
}
