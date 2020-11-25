/*
This program will numerically compute the integral of
                  4/(1+x*x)
from 0 to 1.  The value of this integral is pi -- which
is great since it gives us an easy way to check the answer.
The is the original sequential program.  It uses the timer
from the OpenMP runtime library
History: Written by Tim Mattson, 11/99.
*/
#include <stdio.h>
#include <omp.h>

static long num_steps = 100000000;
double step;

int main ()
{
    int section_threads;
    double pi, sum = 0.0;
    double start_time, run_time;
    double sum_array[4];

    step = 1.0/(double) num_steps;
    section_threads = num_steps / 4;

    start_time = omp_get_wtime();

#pragma omp parallel num_threads(4)
    {
	double tmp;
	int ID = omp_get_thread_num();
	for (int i=section_threads * ID; i< section_threads * (ID+1); i++){
	    tmp = (i-0.5)*step;
	    sum_array[ID] = sum_array[ID] + 4.0/(1.0+tmp*tmp);
	}
    }

    for (int j = 0; j < 4; j++)
	sum = sum + sum_array[j];

    pi = step * sum;
    run_time = omp_get_wtime() - start_time;
    printf("\n pi with %ld steps is %lf in %lf seconds\n ",num_steps,pi,run_time);
}