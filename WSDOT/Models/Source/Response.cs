namespace WSDOT.Models.Source
{
    /// <summary>
    /// The response from the source
    /// </summary>
    public class Response
    {
        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public long CurrentTime { get; set; } = long.MinValue;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public string Description { get; set; } = string.Empty;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public double Distance { get; set; } = double.MinValue;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public long TravelTimeID { get; set; } = long.MinValue;

        /// <inheritdoc />
        public override string ToString() => $"Distance: {this.Distance} miles, Current Time: {this.CurrentTime} minutes";
    }
}