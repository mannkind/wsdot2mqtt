namespace WSDOT.Models.Source
{
    /// <summary>
    /// The response from the source
    /// </summary>
    public record Response
    {
        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public long CurrentTime { get; init; } = long.MinValue;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public string Description { get; init; } = string.Empty;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public double Distance { get; init; } = double.MinValue;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public long TravelTimeID { get; init; } = long.MinValue;
    }
}