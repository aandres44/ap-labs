#include <stdio.h>
#include <stdlib.h>

#define   LOWER  0       /* lower limit of table */
#define   UPPER  300     /* upper limit */
#define   STEP   20      /* step size */

/* print Fahrenheit-Celsius table */

int main(int argc, char *argv[])
{
    int fahr,
        start = atoi(argv[1]),
        end = start,
        increment = 10000;

    if (argc == 4)
    {
        end       = atoi(argv[2]);
        increment = atoi(argv[3]);
    }
    

    for (fahr = start; fahr <= end; fahr = fahr + increment)
	printf("Fahrenheit: %3d, Celcius: %6.1f\n", fahr, (5.0/9.0)*(fahr-32));

    /*for (fahr = LOWER; fahr <= UPPER; fahr = fahr + STEP)
	 *printf("Fahrenheit: %3d, Celcius: %6.1f\n", fahr, (5.0/9.0)*(fahr-32));
     */
    return 0;
    
}
