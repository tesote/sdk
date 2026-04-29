using System.Collections.Generic;
using System.Globalization;

namespace Tesote.Sdk.Internal;

/// <summary>
/// Small helper for building query-parameter dictionaries from typed values.
/// Skips null entries so callers don't pollute the URL with empty params.
/// </summary>
internal sealed class QueryBuilder
{
    private readonly Dictionary<string, string> _values = new();

    public QueryBuilder Add(string key, string? value)
    {
        if (value is null)
        {
            return this;
        }
        _values[key] = value;
        return this;
    }

    public QueryBuilder Add(string key, int? value)
    {
        if (value is null)
        {
            return this;
        }
        _values[key] = value.Value.ToString(CultureInfo.InvariantCulture);
        return this;
    }

    public QueryBuilder Add(string key, long? value)
    {
        if (value is null)
        {
            return this;
        }
        _values[key] = value.Value.ToString(CultureInfo.InvariantCulture);
        return this;
    }

    public QueryBuilder Add(string key, decimal? value)
    {
        if (value is null)
        {
            return this;
        }
        _values[key] = value.Value.ToString(CultureInfo.InvariantCulture);
        return this;
    }

    public QueryBuilder Add(string key, bool? value)
    {
        if (value is null)
        {
            return this;
        }
        _values[key] = value.Value ? "true" : "false";
        return this;
    }

    public IReadOnlyDictionary<string, string>? BuildOrNull()
        => _values.Count == 0 ? null : _values;
}
