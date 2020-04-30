using System.Linq;
using System.Threading.Channels;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using TwoMQTT.Core.DataAccess;
using TwoMQTT.Core.Managers;
using WSDOT.Models.Shared;
using WSDOT.Models.SourceManager;

namespace WSDOT.Managers
{
    /// <summary>
    /// An class representing a managed way to interact with a source.
    /// </summary>
    public class SourceManager : HTTPPollingManager<SlugMapping, FetchResponse, object, Resource, Command>
    {
        /// <summary>
        /// Initializes a new instance of the SourceManager class.
        /// </summary>
        /// <param name="logger"></param>
        /// <param name="sharedOpts"></param>
        /// <param name="opts"></param>
        /// <param name="outgoingData"></param>
        /// <param name="incomingCommand"></param>
        /// <param name="sourceDAO"></param>
        /// <returns></returns>
        public SourceManager(ILogger<SourceManager> logger, IOptions<Models.Shared.Opts> sharedOpts,
            IOptions<Models.SourceManager.Opts> opts, ChannelWriter<Resource> outgoingData,
            ChannelReader<Command> incomingCommand,
            IHTTPSourceDAO<SlugMapping, Command, FetchResponse, object> sourceDAO) :
            base(logger, outgoingData, incomingCommand, sharedOpts.Value.Resources, opts.Value.PollingInterval, sourceDAO)
        {
            this.Opts = opts.Value;
            this.SharedOpts = sharedOpts.Value;
        }

        /// <inheritdoc />
        protected override void LogSettings() =>
            this.Logger.LogInformation(
                $"ApiKey: {this.Opts.ApiKey}\n" +
                $"PollingInterval: {this.Opts.PollingInterval}\n" +
                $"Resources: {string.Join(',', this.SharedOpts.Resources.Select(x => $"{x.TravelTimeID}:{x.Slug}"))}\n" +
                $""
            );

        /// <inheritdoc />
        protected override Resource MapResponse(FetchResponse src) =>
            new Resource
            {
                TravelTimeID = src.TravelTimeID,
                Description = src.Description,
                CurrentTime = src.CurrentTime,
                Distance = src.Distance,
                Closed = src.CurrentTime == 0,
            };

        /// <summary>
        /// The options for the source.
        /// </summary>
        private readonly Models.SourceManager.Opts Opts;

        /// <summary>
        /// The options that are shared.
        /// </summary>
        private readonly Models.Shared.Opts SharedOpts;
    }
}
