using System;
using System.Net.Http;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using Newtonsoft.Json;
using TwoMQTT.Core.DataAccess;
using WSDOT.Models.Shared;

namespace WSDOT.DataAccess
{
    /// <summary>
    /// An class representing a managed way to interact with a source.
    /// </summary>
    public class SourceDAO : SourceDAO<SlugMapping, Command, Models.SourceManager.FetchResponse, object>
    {
        /// <summary>
        /// Initializes a new instance of the SourceDAO class.
        /// </summary>
        /// <param name="logger"></param>
        /// <param name="httpClientFactory"></param>
        /// <param name="apiKey"></param>
        /// <returns></returns>
        public SourceDAO(ILogger<SourceDAO> logger, IHttpClientFactory httpClientFactory, string apiKey) :
            base(logger)
        {
            this.ApiKey = apiKey;
            this.Client = httpClientFactory.CreateClient();
        }

        /// <inheritdoc />
        public override async Task<Models.SourceManager.FetchResponse?> FetchOneAsync(SlugMapping data,
            CancellationToken cancellationToken = default)
        {
            try
            {
                return await this.FetchAsync(data.TravelTimeID, cancellationToken);
            }
            catch (Exception e)
            {
                var msg = e is HttpRequestException ? "Unable to fetch from the WSDOT API" :
                          e is JsonException ? "Unable to deserialize response from the WSDOT API" :
                          "Unable to send to the WSDOT API";
                this.Logger.LogError(msg, e);
                return null;
            }
        }

        /// <summary>
        /// The API Key to access the source.
        /// </summary>
        private readonly string ApiKey;

        /// <summary>
        /// The HTTP client used to access the source.
        /// </summary>
        private readonly HttpClient Client;

        /// <summary>
        /// Fetch one response from the source
        /// </summary>
        /// <param name="timeTravelId"></param>
        /// <param name="cancellationToken"></param>
        /// <returns></returns>
        private async Task<Models.SourceManager.FetchResponse?> FetchAsync(long timeTravelId,
            CancellationToken cancellationToken = default)
        {
            this.Logger.LogDebug($"Started finding {timeTravelId} from WSDOT");
            var baseUrl = "https://www.wsdot.wa.gov/Traffic/api/TravelTimes/TravelTimesREST.svc/GetTravelTimeAsJson";
            var query = $"AccessCode={this.ApiKey}&TravelTimeID={timeTravelId}";
            var resp = await this.Client.GetAsync($"{baseUrl}?{query}", cancellationToken);
            resp.EnsureSuccessStatusCode();
            var content = await resp.Content.ReadAsStringAsync();
            var obj = JsonConvert.DeserializeObject<Models.SourceManager.FetchResponse>(content);
            this.Logger.LogDebug($"Finished finding {timeTravelId} from WSDOT");

            return obj;
        }
    }
}