package cloud.puccini;

import org.snakeyaml.engine.v2.api.Load;
import org.snakeyaml.engine.v2.api.LoadSettings;

import java.util.Map;

public class TOSCA
{
	public static Map<Object, Object> Compile( String url ) throws Exception
	{
		Load load = new SnakeYAML.Load( LoadSettings.builder().build() );
		Map<Object, Object> clout = (Map<Object, Object>) load.loadFromString( _Compile( url ) );
		return clout;
	}

	static
	{
		System.loadLibrary( "puccinijni" );
	}

	public static native String _Compile( String url );
}
