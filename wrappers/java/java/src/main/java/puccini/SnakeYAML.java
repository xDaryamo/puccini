package puccini;

import org.snakeyaml.engine.v2.api.ConstructNode;
import org.snakeyaml.engine.v2.api.DumpSettings;
import org.snakeyaml.engine.v2.api.RepresentToNode;
import org.snakeyaml.engine.v2.exceptions.ConstructorException;
import org.snakeyaml.engine.v2.exceptions.YamlEngineException;
import org.snakeyaml.engine.v2.nodes.Node;
import org.snakeyaml.engine.v2.nodes.ScalarNode;
import org.snakeyaml.engine.v2.nodes.Tag;
import org.snakeyaml.engine.v2.representer.StandardRepresenter;

import java.util.Date;
import java.util.HashMap;
import java.util.Optional;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.time.DateTimeException;
import java.time.Instant;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.ZonedDateTime;
import java.time.ZoneId;
import java.time.format.DateTimeFormatter;
import java.time.format.DateTimeFormatterBuilder;
import java.time.format.DateTimeParseException;
import java.time.temporal.TemporalAccessor;
import java.time.temporal.UnsupportedTemporalTypeException;

/**
 * Adds !!timestamp type support to SnakeYAML engine.
 * <p>
 * See: https://yaml.org/type/timestamp.html
 * </p>
 */
public abstract class SnakeYAML
{
	public static Tag timestampTag = new Tag( Tag.PREFIX + "timestamp" );

	public static HashMap<Tag, ConstructNode> tagConstructors = new HashMap<Tag, ConstructNode>();

	static
	{
		tagConstructors.put( timestampTag, new ConstructTimestampNode() );
	}

	public static class Representer extends StandardRepresenter
	{
		public static DateTimeFormatter zonedDateTimeFormatter = DateTimeFormatter.ofPattern( "yyyy-MM-dd'T'HH:mm:ss.SSSSSSSSSXXX" );

		public static DateTimeFormatter localDateTimeFormatter = DateTimeFormatter.ofPattern( "yyyy-MM-dd'T'HH:mm:ss.SSSSSSSSS'Z'" );

		public static DateTimeFormatter localDateFormatter = DateTimeFormatter.ofPattern( "yyyy-MM-dd'T00:00:00Z'" );

		public Representer( DumpSettings settings )
		{
			super( settings );
			RepresentTimestamp representer = new RepresentTimestamp();
			this.representers.put( ZonedDateTime.class, representer );
			this.representers.put( LocalDateTime.class, representer );
			this.representers.put( LocalDate.class, representer );
			this.representers.put( Instant.class, representer );
			this.representers.put( Date.class, representer );
		}

		class RepresentTimestamp implements RepresentToNode
		{
			public Node representData( Object data )
			{
				TemporalAccessor temporalAccessor;

				if( data instanceof Date )
					data = ( (Date) data ).toInstant();

				if( data instanceof Instant )
					// To ZonedDateTime
					temporalAccessor = ( (Instant) data ).atZone( ZoneId.systemDefault() );
				else if( data instanceof TemporalAccessor )
					temporalAccessor = (TemporalAccessor) data;
				else
					// This should never happen
					throw new YamlEngineException( "incompatible timestamp: " + data );

				String value;

				try
				{
					// Has date, time, and timezone
					value = zonedDateTimeFormatter.format( temporalAccessor );
				}
				catch( UnsupportedTemporalTypeException x )
				{
					try
					{
						// Has date and time, but no timezone
						value = localDateTimeFormatter.format( temporalAccessor );
					}
					catch( UnsupportedTemporalTypeException xx )
					{
						try
						{
							// Only date, no time and timezone
							value = localDateFormatter.format( temporalAccessor );
						}
						catch( UnsupportedTemporalTypeException xxx )
						{
							// This should never happen
							throw new YamlEngineException( "cannot format timestamp: " + temporalAccessor );
						}
					}
				}

				return representScalar( timestampTag, value );
			}
		}
	}

	static class ConstructTimestampNode implements ConstructNode
	{
		public static Pattern patternShort = Pattern.compile( "^(?<year>[0-9][0-9][0-9][0-9])-(?<month>[0-9][0-9])-(?<day>[0-9][0-9])$" );

		public static Pattern patternLong = Pattern.compile(
			"^(?<year>[0-9][0-9][0-9][0-9])-(?<month>[0-9][0-9])-(?<day>[0-9][0-9])(?:[Tt]|[ \\t]+)(?<hour>[0-9][0-9]?):(?<minute>[0-9][0-9]):(?<second>[0-9][0-9])(?:(?<fraction>\\.[0-9]*))?(?:(?:[ \\t]*)(?:Z|(?<tzhour>[-+][0-9][0-9]?)(?::(?<tzminute>[0-9][0-9]))?))?$" );

		public static DateTimeFormatter formatter = new DateTimeFormatterBuilder().parseLenient().appendPattern( "yyyy-M-d[['T'H:m:s[.S]][XXX]]" ).toFormatter();

		public Object construct( Node node )
		{
			String value = ( (ScalarNode) node ).getValue();

			// Validate
			Matcher matcher = patternShort.matcher( value );
			if( !matcher.find() )
			{
				matcher = patternLong.matcher( value );
				if( !matcher.find() )
					throw new ConstructorException( null, Optional.empty(), "malformed timestamp", node.getStartMark() );

				// Reformat to canonical form (without whitespace)
				StringBuilder s = new StringBuilder();
				s.append( String.format( "%s-%s-%sT%s:%s:%s%s", matcher.group( "year" ), matcher.group( "month" ), matcher.group( "day" ), matcher.group( "hour" ), matcher.group( "minute" ), matcher.group( "second" ),
					matcher.group( "fraction" ) ) );
				String tzhour = matcher.group( "tzhour" );
				String tzminute = matcher.group( "tzminute" );
				if( ( tzhour != "" ) && ( tzminute != "" ) )
					s.append( String.format( "%s:%s", tzhour, tzminute ) );
				else
					s.append( 'Z' );
				value = s.toString();
			}

			try
			{
				TemporalAccessor temporalAccessor = formatter.parse( value );

				try
				{
					// Has date, time, and timezone
					return ZonedDateTime.from( temporalAccessor );
				}
				catch( DateTimeException x )
				{
					try
					{
						// Has date and time, but no timezone
						return LocalDateTime.from( temporalAccessor );
					}
					catch( DateTimeException xx )
					{
						// Only date, no time and timezone
						return LocalDate.from( temporalAccessor );
					}
				}
			}
			catch( DateTimeParseException x )
			{
				throw new ConstructorException( null, Optional.empty(), "cannot parse timestamp: " + value, node.getStartMark() );
			}
		}
	}
}