namespace WSDOT.Models.Shared
{
    /// <summary>
    /// The shared resource across the application
    /// </summary>
    public class Resource
    {
        public long TravelTimeID { get; set; } = 0;
        public string Description { get; set; } = string.Empty;
        public long CurrentTime { get; set; } = 0;
        public double Distance { get; set; } = 0.0;
        public bool Closed { get; set; } = false;

        public override string ToString() => $"Description: {this.Description}, Distance: {this.Distance} miles, Current Time: {this.CurrentTime} minutes";

        public static Resource From(SourceManager.Response obj) => 
            new Resource
            {
                TravelTimeID = obj.TravelTimeID,
                Description = obj.Description,
                CurrentTime = obj.CurrentTime,
                Distance = obj.Distance,
                Closed = obj.CurrentTime == 0,
            };
    }
}
