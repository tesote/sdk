using System;

namespace Tesote.Sdk.Errors;

/// <summary>Bad SDK configuration; raised at construction.</summary>
public sealed class ConfigException : TesoteException
{
    /// <summary>Construct with a message only.</summary>
    public ConfigException(string message)
        : base(message, "CONFIG", 0, null, null, null, null, null, 0, null)
    {
    }

    /// <summary>Construct with a message and underlying cause.</summary>
    public ConfigException(string message, Exception cause)
        : base(message, "CONFIG", 0, null, null, null, null, null, 0, cause)
    {
    }
}
