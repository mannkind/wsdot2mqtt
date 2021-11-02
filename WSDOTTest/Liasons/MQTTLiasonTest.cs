using System.Linq;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using Moq;
using TwoMQTT.Utils;
using WSDOT.Liasons;
using WSDOT.Models.Options;
using WSDOT.Models.Shared;

namespace WSDOTTest.Liasons;

[TestClass]
public class MQTTLiasonTest
{
    [TestMethod]
    public void MapDataTest()
    {
        var tests = new[] {
                new {
                    Q = new SlugMapping { TravelTimeID = BasicTravelTimeID, Slug = BasicSlug },
                    Resource = new Resource { TravelTimeID = BasicTravelTimeID, CurrentTime = BasicTime },
                    Expected = new { TravelTimeID = BasicTravelTimeID, CurrentTime = BasicTime.ToString(), Slug = BasicSlug, Found = true }
                },
                new {
                    Q = new SlugMapping { TravelTimeID = BasicTravelTimeID, Slug = BasicSlug },
                    Resource = new Resource { TravelTimeID = BasicTravelTimeID-85 },
                    Expected = new { TravelTimeID = 0L, CurrentTime = string.Empty, Slug = string.Empty, Found = false }
                },
            };

        foreach (var test in tests)
        {
            var logger = new Mock<ILogger<MQTTLiason>>();
            var generator = new Mock<IMQTTGenerator>();
            var sharedOpts = Options.Create(new SharedOpts
            {
                Resources = new[] { test.Q }.ToList(),
            });

            generator.Setup(x => x.BuildDiscovery(It.IsAny<string>(), It.IsAny<string>(), It.IsAny<System.Reflection.AssemblyName>(), false))
                .Returns(new TwoMQTT.Models.MQTTDiscovery());
            generator.Setup(x => x.StateTopic(test.Q.Slug, It.IsAny<string>()))
                .Returns($"totes/{test.Q.Slug}/topic/{nameof(Resource.CurrentTime)}");

            var mqttLiason = new MQTTLiason(logger.Object, generator.Object, sharedOpts);
            var results = mqttLiason.MapData(test.Resource);
            var actual = results.FirstOrDefault();

            Assert.AreEqual(test.Expected.Found, results.Any(), "The mapping should exist if found.");
            if (test.Expected.Found)
            {
                Assert.IsTrue(actual.topic.Contains(test.Expected.Slug), "The topic should contain the expected TravelTimeID.");
                Assert.AreEqual(test.Expected.CurrentTime, actual.payload, "The payload be the expected CurrentTime.");
            }
        }
    }

    [TestMethod]
    public void DiscoveriesTest()
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
            var logger = new Mock<ILogger<MQTTLiason>>();
            var generator = new Mock<IMQTTGenerator>();
            var sharedOpts = Options.Create(new SharedOpts
            {
                Resources = new[] { test.Q }.ToList(),
            });

            generator.Setup(x => x.BuildDiscovery(test.Q.Slug, It.IsAny<string>(), It.IsAny<System.Reflection.AssemblyName>(), false))
                .Returns(new TwoMQTT.Models.MQTTDiscovery());

            var mqttLiason = new MQTTLiason(logger.Object, generator.Object, sharedOpts);
            var results = mqttLiason.Discoveries();
            var result = results.FirstOrDefault();

            Assert.IsNotNull(result, "A discovery should exist.");
        }
    }

    private static string BasicSlug = "totallyaslug";
    private static long BasicTime = 52;
    private static long BasicTravelTimeID = 15873525;
}
