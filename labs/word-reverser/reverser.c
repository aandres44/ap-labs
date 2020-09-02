#include <stdio.h>


#define IN   1   /* inside a word */
#define OUT  0   /* outside a word */


int main(){
    // Place your magic here

    //string
    int i = 0,
        j = 0;
    char c[20];
    while ((c[i] = getchar()) != EOF)
    {
        i++;
        j++;
    }
    
    i = j - 1;

    while (i >= 0)
    {
        printf("%c", c[i]);
        i--;
    }
    

    return 0;
}
