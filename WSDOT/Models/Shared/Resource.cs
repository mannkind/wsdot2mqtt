namespace WSDOT.Models.Shared
{
    /// <summary>
    /// The shared resource across the application
    /// </summary>
    public record Resource
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
        public string Description { get; init; } = string.Empty;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public long CurrentTime { get; init; } = 0;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public double Distance { get; init; } = 0.0;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public bool Closed { get; init; } = false;

        /// <inheritdoc />
        public override string ToString() => $"Description: {this.Description}, Distance: {this.Distance} miles, Current Time: {this.CurrentTime} minutes";
    }
}
