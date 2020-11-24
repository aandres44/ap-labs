#include <stdio.h>
#include <string.h>
#include <syslog.h>
#include <stdarg.h>
#include "logger.h"

#define RESET		0
#define BRIGHT 		1



#define RED		    1
#define GREEN		2
#define YELLOW		3
#define BLUE		4
#define MAGENTA		5

#define	WHITE		7
#define BLACK 		10

int logtyp = 0;

int initLogger(char *logType) {
    if (strcmp(logType, "") == 0 || strcmp(logType, "stdout") == 0) {
        logtyp = 0;
    } else if (strcmp(logType, "syslog") == 0){
        logtyp = 1;
    } else {
        perror("Please specify a type <stdout> or <syslog>");
        return 1;
    }
    printf("Initializing Logger on: %s\n", logType);
    return 0;
}

int infof(const char *format, ...) {
    va_list args;
    va_start(args,format);
    textcolor(BRIGHT, BLUE, BLACK);
    if(logtyp){
        openlog("INFO", LOG_CONS | LOG_PID | LOG_NDELAY, LOG_LOCAL1);
        vsyslog(LOG_INFO, format, args);
    }else{
        vfprintf(stdout,format,args);
    }

    va_end(args);
    if (logtyp == 0) {
		fprintf(stdout, "\n");
		fflush(stdout);
	}
    textcolor(RESET, WHITE, BLACK);
    return 0;
}

int warnf(const char *format, ...) {
    va_list args;
    va_start(args,format);
    textcolor(BRIGHT, YELLOW, BLACK);

    if (logtyp){
        openlog("WARN", LOG_CONS | LOG_PID | LOG_NDELAY, LOG_LOCAL1);
		vsyslog(LOG_WARNING, format, args);
		closelog();
    }
	else {
        vfprintf(stdout, format, args);
	}

    va_end(args);
    if (logtyp == 0) {
		fprintf(stdout, "\n");
		fflush(stdout);
	}
    textcolor(RESET, WHITE, BLACK);
    return 0;
}

int errorf(const char *format, ...) {

    va_list args;
    va_start(args,format);

    textcolor(BRIGHT, RED, BLACK);

    if (logtyp){
        openlog("ERROR", LOG_CONS | LOG_PID | LOG_NDELAY, LOG_LOCAL1);
		vsyslog(LOG_ERR, format, args);
		closelog();
    }
	else {
        vfprintf(stdout, format, args);
	}

    va_end(args);
    if (logtyp == 0) {
		fprintf(stdout, "\n");
		fflush(stdout);
	}
    textcolor(RESET, WHITE, BLACK);
    return 0;
}
int panicf(const char *format, ...) {
    va_list args;
    va_start(args,format);
    textcolor(BRIGHT, MAGENTA, BLACK);

    if (logtyp){
        openlog("PANIC", LOG_CONS | LOG_PID | LOG_NDELAY, LOG_LOCAL1);
		vsyslog(LOG_EMERG, format, args);
		closelog();
    }
	else {
        vfprintf(stdout, format, args);
	}
    
    va_end(args);
    if (logtyp == 0) {
		fprintf(stdout, "\n");
		fflush(stdout);
	}
    textcolor(RESET, WHITE, BLACK);
    return 0;
}

void textcolor(int attr, int fg, int bg)
{	
    char command[13];
	sprintf(command, "%c[%d;%d;%dm", 0x1B, attr, fg + 30, bg + 40);
	printf("%s", command);
}