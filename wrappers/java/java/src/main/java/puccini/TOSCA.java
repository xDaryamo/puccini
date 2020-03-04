package puccini;

import org.snakeyaml.engine.v2.api.Load;
import org.snakeyaml.engine.v2.api.LoadSettings;

import java.util.Map;

public class TOSCA
{
	public static Map<Object, Object> Compile( String url ) throws Exception
	{
		LoadSettings settings = LoadSettings.builder().setTagConstructors( SnakeYAML.tagConstructors ).build();
		Load load = new Load( settings );
		Map<Object, Object> clout = (Map<Object, Object>) load.loadFromString( _Compile( url ) );
		return clout;
	}

	static
	{
		System.loadLibrary( "puccinijni" );
	}

	public static native String _Compile( String url );
}
