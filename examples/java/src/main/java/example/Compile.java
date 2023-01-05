package example;

import org.snakeyaml.engine.v2.api.Dump;
import org.snakeyaml.engine.v2.api.DumpSettings;
import org.snakeyaml.engine.v2.api.Load;
import org.snakeyaml.engine.v2.api.LoadSettings;
import cloud.puccini.Problems;
import cloud.puccini.SnakeYAML;
import cloud.puccini.TOSCA;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class Compile
{
	public static void main( String[] args )
	{
		if( args.length >= 1 )
		{
			String url = "";
			Map<String, Object> inputs = new HashMap<String, Object>();
			List<String> quirks = new ArrayList<String>();
			boolean resolve = true;
			boolean coerce = true;

			Load load = new SnakeYAML.Load( LoadSettings.builder().build() );
			for ( String arg: args )
			{
				if ( arg.startsWith( "--input=" ) ) {
					String[] s = arg.substring( 8 ).split( "=" );
					inputs.put( s[0], load.loadFromString( s[1] ) );
				} else if ( arg.startsWith( "--quirk=" ) ) {
					quirks.add(arg.substring( 8 ));
				} else if ( arg.startsWith( "--resolve=" ) ) {
					resolve = arg.substring( 10 ) == "true";
				} else if ( arg.startsWith( "--coerce=" ) ) {
					coerce = arg.substring( 9 ) == "true";
				} else {
					url = arg;
				}
			}

			Dump dump = new SnakeYAML.Dump( DumpSettings.builder().build() );
			try
			{
				Object clout = TOSCA.Compile( url, inputs, quirks, resolve, coerce );
				System.out.print( dump.dumpToString( clout ) );
			}
			catch( Problems x )
			{
				System.err.println( "Problems:" );
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