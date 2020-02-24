namespace WSDOT.Models.Shared
{
    /// <summary>
    /// The shared key info => slug mapping across the application
    /// </summary>
    public class SlugMapping
    {
        public long TravelTimeID { get; set; } = 0;
        public string Slug { get; set; } = string.Empty;

        public override string ToString() => $"Time Travel ID: {this.TravelTimeID}, Slug: {this.Slug}";
    }
}
