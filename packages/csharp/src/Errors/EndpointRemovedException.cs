namespace Tesote.Sdk.Errors;

/// <summary>
/// Raised when calling a method whose upstream endpoint is gone in this API
/// version. The SDK keeps the method per the back-compat policy but throws this
/// pointing at the replacement.
/// </summary>
public sealed class EndpointRemovedException : TesoteException
{
    /// <summary>Construct with a human-readable message.</summary>
    public EndpointRemovedException(string message)
        : base(message, "ENDPOINT_REMOVED", 0, null, null, null, null, null, 0, null)
    {
    }
}
