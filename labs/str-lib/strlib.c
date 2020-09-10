#include <stdlib.h>
#include <stdio.h>



int mystrlen(char *str) {

    int size = 0;
    while(str[size] != '\0') {
        size++;
    }
    return size;
}



char *mystradd(char *origin, char *addition) {

    int len = mystrlen(origin);
    int i;
    for(i = 0; addition[i] != '\0'; i++){
        origin[len + i] = addition[i];
    }
    origin[len + i] = '\0';
    return origin;

}

int mystrfind(char *origin, char *substr) {

    int len = mystrlen(origin);
    int sublen = mystrlen(substr);
    int i, j;

    for(i = 0; i < len; i++){

        if(origin[i] == substr[0]){

            for (j = 1; j < sublen; j++){

                if (origin[i + j] != substr[j]){
                    break;
                }
            }
            if(j == sublen){
                return i;
            }
        }
    }

    return -1;

    
}
