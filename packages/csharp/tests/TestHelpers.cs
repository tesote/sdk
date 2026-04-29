using System;
using Tesote.Sdk;
using Tesote.Sdk.V1;
using Tesote.Sdk.V2;

namespace Tesote.Sdk.Tests;

internal static class TestHelpers
{
    public static V1Client NewV1(string baseUrl)
    {
        return new V1Client(new ClientOptions
        {
            ApiKey = "sk_test_abcd1234",
            BaseUrl = baseUrl,
            RequestTimeout = TimeSpan.FromSeconds(2),
            RetryPolicy = new RetryPolicy(1, TimeSpan.FromMilliseconds(1), TimeSpan.FromMilliseconds(2), false),
        });
    }

    public static V2Client NewV2(string baseUrl)
    {
        return new V2Client(new ClientOptions
        {
            ApiKey = "sk_test_abcd1234",
            BaseUrl = baseUrl,
            RequestTimeout = TimeSpan.FromSeconds(2),
            RetryPolicy = new RetryPolicy(1, TimeSpan.FromMilliseconds(1), TimeSpan.FromMilliseconds(2), false),
        });
    }
}
