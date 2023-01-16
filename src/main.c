#include <stdio.h>
#include <stdlib.h>
#include <sys/dir.h>
#include <sys/stat.h>
#include <time.h>

int main(int argc, char *argv[]) {
    if (argc < 2) {
        printf("Usage: %s <directory>", argv[0]);
        exit(1);
    }

    DIR *dir = opendir(argv[1]);
    if (dir == NULL) {
        perror("opendir");
        exit(1);
    }

    struct dirent *entry;

    while ((entry = readdir(dir)) != NULL) {
        struct stat st;
        printf("%s\n", entry->d_name);
    }
}