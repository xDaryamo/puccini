package cloud.puccini;

import java.util.List;

public class Problems extends Exception
{
    public List<Object> problems;

    public Problems( List<Object> problems )
    {
        super( "problems" );
        this.problems = problems;
    }
}