package example;

import java.util.Map;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;
import puccini.TOSCA;

public class Compile
{
	public static void main( String[] args )
	{
		if( args.length >= 1 )
		{
			try
			{
				Map<Object, Object> clout = TOSCA.Compile( args[0] );
				ObjectMapper mapper = new ObjectMapper( new YAMLFactory() );
				mapper.writeValue( System.out, clout );
			}
			catch( Exception x )
			{
				System.err.println( x );
				System.exit( 1 );
			}
		}
		else
		{
			System.err.println( "no URL provided" );
			System.exit( 1 );
		}
	}
}