#include <cassert>
#include "util.h"

using namespace std;

int main() {
    printf("string_hash test\n");
    assert(string_hash("0") == 48 % HASH_MOD);
    assert(string_hash("A") == 65 % HASH_MOD);
    assert(string_hash("http") == 3213448 % HASH_MOD);
    printf("ok\n");

    printf("hash test\n");
    hash_add(0, "https://www.google.com/");
    hash_add(1, "https://ecl.me/");
    hash_add(1, "https://github.com/");
    assert(strcmp(fgeturl(hash_get(0)), "https://www.google.com/") == 0);
    auto hash_file1 = hash_get(1);
    assert(strcmp(fgeturl(hash_file1), "https://ecl.me/") == 0);
    assert(strcmp(fgeturl(hash_file1), "https://github.com/") == 0);
    remove_tmp();
    printf("ok\n");
    return 0;
}
