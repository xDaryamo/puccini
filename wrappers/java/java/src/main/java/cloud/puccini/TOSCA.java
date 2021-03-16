package cloud.puccini;

import org.snakeyaml.engine.v2.api.Dump;
import org.snakeyaml.engine.v2.api.DumpSettings;
import org.snakeyaml.engine.v2.api.Load;
import org.snakeyaml.engine.v2.api.LoadSettings;

import java.util.List;
import java.util.Map;

public class TOSCA
{
	public static Object Compile( String url, Map<String, Object> inputs ) throws Exception
	{
		Load load = new SnakeYAML.Load( LoadSettings.builder().build() );
		Dump dump = new SnakeYAML.Dump( DumpSettings.builder().build() );

		String inputs_ = dump.dumpToString( inputs );
		Map<Object, Object> result = (Map<Object, Object>) load.loadFromString( _Compile( url, inputs_ ) );

		if ( result.containsKey( "problems" ) )
		{
			Object problems = result.get( "problems" );
			if ( problems instanceof List<?> )
			{
				throw new Problems( (List<Object>) problems );
			}
		}
		else if ( result.containsKey( "error" ) )
		{
			throw new Exception( result.get( "error" ).toString() );
		}
		else
		{
			return result.get( "clout" );
		}

		return result;
	}

	static
	{
		System.loadLibrary( "puccinijni" );
	}

	public static native String _Compile( String url, String inputs );
}
