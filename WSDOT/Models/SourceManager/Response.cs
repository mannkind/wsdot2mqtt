namespace WSDOT.Models.SourceManager
{
    /// <summary>
    /// The response from the source
    /// </summary>
    public class Response
    {
        public long AverageTime { get; set; } = 0;
        public long CurrentTime { get; set; } = 0;
        public string Description { get; set; } = string.Empty;
        public double Distance { get; set; } = 0.0;
        public string Name { get; set; } = string.Empty;
        public string TimeUpdated { get; set; } = string.Empty;
        public long TravelTimeID { get; set; } = 0;
        public bool Ok { get; set; } = false;

        public override string ToString() => $"Distance: {this.Distance} miles, Current Time: {this.CurrentTime} minutes";
    }
}