using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using Moq;
using WSDOT.DataAccess;
using WSDOT.Liasons;
using WSDOT.Models.Options;
using WSDOT.Models.Shared;

namespace WSDOTTest.Liasons;

[TestClass]
public class SourceLiasonTest
{
    [TestMethod]
    public async Task FetchAllAsyncTest()
    {
        var tests = new[] {
                new {
                    Q = new SlugMapping { TravelTimeID = BasicTravelTimeID, Slug = BasicSlug },
                    Resource = new Resource { TravelTimeID = BasicTravelTimeID, CurrentTime = BasicTime },
                    Expected = new { TravelTimeID = BasicTravelTimeID, CurrentTime = BasicTime, Slug = BasicSlug }
                },
            };

        foreach (var test in tests)
        {
            var logger = new Mock<ILogger<SourceLiason>>();
            var sourceDAO = new Mock<ISourceDAO>();
            var opts = Options.Create(new SourceOpts());
            var sharedOpts = Options.Create(new SharedOpts
            {
                Resources = new[] { test.Q }.ToList(),
            });

            sourceDAO.Setup(x => x.FetchOneAsync(test.Q, It.IsAny<CancellationToken>()))
                 .ReturnsAsync(new WSDOT.Models.Source.Response
                 {
                     TravelTimeID = test.Expected.TravelTimeID,
                     CurrentTime = test.Expected.CurrentTime,
                 });

            var sourceLiason = new SourceLiason(logger.Object, sourceDAO.Object, opts, sharedOpts);
            await foreach (var result in sourceLiason.ReceiveDataAsync())
            {
                Assert.AreEqual(test.Expected.TravelTimeID, result.TravelTimeID);
                Assert.AreEqual(test.Expected.CurrentTime, result.CurrentTime);
            }
        }
    }

    private static string BasicSlug = "totallyaslug";
    private static long BasicTime = 52;
    private static long BasicTravelTimeID = 15873525;
}
