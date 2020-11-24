#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <pthread.h>
#include "logger.h"
#define _NUM_THREADS_ 2000

long **buffers;
int NUM_BUFFERS = 0;
char *RESULT_MATRIX_FILE;
pthread_mutex_t *mtx;
pthread_t threads[_NUM_THREADS_];

struct args_struct
{
    int row;
    int col;
    long *matA;
    long *matB;
};

long * readMatrix(char *filename);
long * multiply(long *matA, long *matB);
long * getColumn(int col, long *matrix);
long * getRow(int row, long *matrix);
long dotProduct(long *vecA, long *vecB);
int saveResult(long *result);
int getLock();
int releaseLock(int lock);

int main(int argc, char **argv)
{
    initLogger("stdout");

    if (argc < 5)
    {
        errorf("Error: not enough arguments provided.\n");
        errorf("Usage: ./multiplier -n <buffers> -out <file>\n");
        return -1;
    }

    if (strcmp(argv[1], "-n") == 0 && strcmp(argv[3], "-out") == 0)
    {
        NUM_BUFFERS = atoi(argv[2]);
        RESULT_MATRIX_FILE = argv[4];
    }
    else if (strcmp(argv[1], "-out") == 0 && strcmp(argv[3], "-n") == 0)
    {
        RESULT_MATRIX_FILE = argv[2];
        NUM_BUFFERS = atoi(argv[4]);
    }else{
        errorf("Error: Wrong number of arguments.\n");
        errorf("Usage: ./multiplier -n <buffers> -out <file>\n");
        return -1;
    }

    buffers = malloc(NUM_BUFFERS * sizeof(long *));
    mtx = malloc(NUM_BUFFERS * sizeof(pthread_mutex_t));


    for(int i = 0; i < NUM_BUFFERS; i++){
        pthread_mutex_init(&mtx[i], NULL);
    }

    long *matA = readMatrix("matA.dat");

    long *matB = readMatrix("matB.dat");

    

    long *res = multiply(matA, matB);

    infof("Operation finished\n");

    saveResult(res);
    free(matA);
    free(matB);
    free(mtx);
    free(buffers);
    free(res);

    return 0;
}

long * getRow(int row, long *matrix){
    long *r = malloc(_NUM_THREADS_ * sizeof(long));
    int rowStart = _NUM_THREADS_ * (row);
    for(int i = 0;i < _NUM_THREADS_; i++){
        r[i] = matrix[rowStart + i];
    }
    return r;
}

long * getColumn(int column, long *matrix){
    long *c = malloc(_NUM_THREADS_ * sizeof(long));
    for(int i = 0; i < _NUM_THREADS_; i++){
        c[i] = matrix[column];
        column += _NUM_THREADS_;
    }
    return c;
}

int getLock() {
    for(int i = 0; i < NUM_BUFFERS; i++){
        if(pthread_mutex_trylock(&mtx[i]) == 0){
            return i;
        }   
    }
    return -1;
}

int releaseLock(int lock) {
    return pthread_mutex_unlock(&mtx[lock]);
}

long * readMatrix( char *filename){
    int size = 0;
    FILE *f = fopen(filename,"r");
    if(f == NULL){
        errorf("Error: can't open Matrix data file");
        exit(2);
    }
    char c;
    while((c = fgetc(f)) != EOF){
        if(c == '\n')
            size++;
    }
    rewind(f);
    long *matrix = malloc(size * sizeof(long));
    int index = 0;
    while(fscanf(f,"%ld",&matrix[index]) != EOF){
        index++;
    }
    fclose(f);
    infof("Task: Read matrix, ended succesfully");
    return matrix;   
}

long dotProduct(long *vecA, long *vecB) {
    long sum = 0;
    for(int i = 0; i < 2000; i ++) {
        sum += (vecA[i] * vecB[i]);
    }
    return sum;
}

long performMult(struct args_struct *data){
    int bufA = -1, bufB = -1;
    while(bufA == -1 || bufB == -1){
        if(bufA == -1){
            bufA = getLock();
        }
        if(bufB == -1){
            bufB = getLock();
        }
    }
    buffers[bufA] = getRow(data->row, data->matA);
    buffers[bufB] = getColumn(data->col, data->matB);
    long res = dotProduct(buffers[bufA], buffers[bufB]);
    free(buffers[bufA]);
    free(buffers[bufB]);
    free(data);
    releaseLock(bufA);
    releaseLock(bufB);
    return res;
}

long * multiply(long *matA, long *matB){
    infof("Multiplying matrices task starting...");
    long *result = malloc(_NUM_THREADS_*_NUM_THREADS_ * sizeof(long));
    for(int i = 0; i < _NUM_THREADS_; i++){
        for(int j = 0;j < _NUM_THREADS_; j++){
            struct args_struct *ar = malloc(sizeof(struct args_struct));
            ar -> row = i;
            ar -> col = j;
            ar -> matA = matA;
            ar -> matB = matB;
            pthread_create(&threads[j],NULL,(void * (*) (void *))performMult,(void *)ar);
        }
        for(int j = 0; j < _NUM_THREADS_; j++){
            void *res;
            pthread_join(threads[j], &res);
            result[_NUM_THREADS_ * j + i] = (long) res;
        }
    }
    infof("Multiplying matrices task ended succesfully");
    return result;
}

int saveResult(long *result) {    
    FILE *file;
    file = fopen(RESULT_MATRIX_FILE, "w+");
    if(file==NULL){
        return -1;
    }
    int size=_NUM_THREADS_*_NUM_THREADS_;
    for(int i = 0; i < size; i++) {
        fprintf(file, "%ld\n", result[i]);
    }
    fclose(file);
    return 0;
}