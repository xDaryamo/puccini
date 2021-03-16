package example;

import org.snakeyaml.engine.v2.api.Dump;
import org.snakeyaml.engine.v2.api.DumpSettings;
import cloud.puccini.Problems;
import cloud.puccini.SnakeYAML;
import cloud.puccini.TOSCA;

import java.util.Map;
import java.util.HashMap;

public class Compile
{
	public static void main( String[] args )
	{
		if( args.length >= 1 )
		{
			try
			{
				Map<String, Object> inputs = new HashMap<String, Object>();
				Object clout = TOSCA.Compile( args[0], inputs );
				Dump dump = new SnakeYAML.Dump( DumpSettings.builder().build() );
				System.out.print( dump.dumpToString( clout ) );
			}
			catch( Problems x )
			{
				System.err.println( "Problems:" );
				Dump dump = new SnakeYAML.Dump( DumpSettings.builder().build() );
				for ( Object problem : x.problems )
				{
					System.err.print( dump.dumpToString( problem ) );
				}
				System.exit( 1 );
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