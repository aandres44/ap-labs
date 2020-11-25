#include <stdio.h>
#include <omp.h>
int main() {
    #pragma  omp parallel num_threads(8)
    {
        printf("Hello World from thread number %d\n",omp_get_thread_num());
    }
    return 0;
}