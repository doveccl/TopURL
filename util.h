#ifndef UTIL_H
#define UTIL_H

#include <cstdio>
#include <string>

#define HASH_MOD 131
#define MAX_URL_LEN 8192

void hash_add(int, const char *);
FILE *hash_get(int no);
const char *fgeturl(FILE *);
void remove_tmp();
int string_hash(const char *);

#endif // UTIL_H
