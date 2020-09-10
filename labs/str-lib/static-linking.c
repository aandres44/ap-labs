#include <stdio.h>
#include <string.h>

int mystrlen(char *);
char *mystradd(char *, char *);
int mystrfind(char *, char *);

int main(int argc, char *argv[]) {

    if(argc > 3) {
        if(strcmp(argv[1], "-add") == 0) {
            int len = mystrlen(argv[2]);
            char* str = mystradd(argv[2], argv[3]);
            int newLen = mystrlen(str);
            printf("Initial Length      : %i\n", len);
            printf("New String          : %s\n", str);
            printf("New Length          : %i\n", newLen);
        } else if(strcmp(argv[1], "-find") == 0){
            int pos = mystrfind(argv[2], argv[3]);
            printf("['%s'] sring was found at [%i] position\n", argv[3], pos);
        } else{
            printf("Invalid argument\n");
        }
    }

    return 0;
}
