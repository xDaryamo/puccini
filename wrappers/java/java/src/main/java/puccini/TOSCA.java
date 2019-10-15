package puccini;

import java.util.Map;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;

public class TOSCA
{
	static
	{
		System.loadLibrary( "puccinijni" );
	}

	public static native String _Compile( String url );

	public static Map<Object, Object> Compile( String url ) throws Exception
	{
		ObjectMapper mapper = new ObjectMapper( new YAMLFactory() );
		try
		{
			@SuppressWarnings("unchecked")
			Map<Object, Object> clout = mapper.readValue( _Compile( url ), Map.class );
			return clout;
		}
		catch( JsonProcessingException x )
		{
			throw x;
		}
	}
}
