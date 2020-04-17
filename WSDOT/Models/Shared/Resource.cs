namespace WSDOT.Models.Shared
{
    /// <summary>
    /// The shared resource across the application
    /// </summary>
    public class Resource
    {
        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public long TravelTimeID { get; set; } = 0;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public string Description { get; set; } = string.Empty;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public long CurrentTime { get; set; } = 0;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public double Distance { get; set; } = 0.0;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public bool Closed { get; set; } = false;

        /// <inheritdoc />
        public override string ToString() => $"Description: {this.Description}, Distance: {this.Distance} miles, Current Time: {this.CurrentTime} minutes";
    }
}
