using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using TwoMQTT.Interfaces;
using TwoMQTT.Liasons;
using WSDOT.DataAccess;
using WSDOT.Models.Options;
using WSDOT.Models.Shared;
using WSDOT.Models.Source;

namespace WSDOT.Liasons;

/// <summary>
/// A class representing a managed way to interact with a source.
/// </summary>
public class SourceLiason : PollingSourceLiasonBase<Resource, SlugMapping, ISourceDAO, SharedOpts>, ISourceLiason<Resource, object>
{
    public SourceLiason(ILogger<SourceLiason> logger, ISourceDAO sourceDAO,
        IOptions<SourceOpts> opts, IOptions<SharedOpts> sharedOpts) :
        base(logger, sourceDAO, sharedOpts)
    {
        this.Logger.LogInformation(
            "ApiKey: {apiKey}\n" +
            "PollingInterval: {pollingInterval}\n" +
            "Resources: {@resources}\n" +
            "",
            opts.Value.ApiKey,
            opts.Value.PollingInterval,
            sharedOpts.Value.Resources
        );
    }

    /// <inheritdoc />
    protected override async Task<Resource?> FetchOneAsync(SlugMapping key, CancellationToken cancellationToken)
    {
        var result = await this.SourceDAO.FetchOneAsync(key, cancellationToken);
        return result switch
        {
            Response => new Resource
            {
                TravelTimeID = result.TravelTimeID,
                Description = result.Description,
                CurrentTime = result.CurrentTime,
                Distance = result.Distance,
                Closed = result.CurrentTime == 0,
            },
            _ => null,
        };
    }
}
