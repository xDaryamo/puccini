package example;

import org.snakeyaml.engine.v2.api.Dump;
import org.snakeyaml.engine.v2.api.DumpSettings;
import cloud.puccini.SnakeYAML;
import cloud.puccini.TOSCA;

import java.util.Map;

public class Compile
{
	public static void main( String[] args )
	{
		if( args.length >= 1 )
		{
			try
			{
				Map<Object, Object> clout = TOSCA.Compile( args[0] );
				Dump dump = new SnakeYAML.Dump( DumpSettings.builder().build() );
				System.out.print( dump.dumpToString( clout ) );
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