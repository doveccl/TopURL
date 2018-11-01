#include "util.h"

using namespace std;

FILE *tmp[HASH_MOD];

/**
 * 根据哈希值打开对应小文件
 * @param 哈希值
 * @param 文件读写模式
 * @return 文件指针
 */
FILE *hash_open(int no, const char *mode) {
    char buf[20];
    sprintf(buf, "HASH_%d", no);
    return fopen(buf, mode);
}

/**
 * 判断字符是否为 url 合法字符
 * @param 待检测字符
 * @return 检测结果
 */
bool url_char(char c) {
    return ' ' < c && c <= '~';
}

/**
 * 将一个 url 加入对应的哈希小文件中
 * @param 哈希值
 * @param url url内容
 */
void hash_add(int no, const char *url) {
    if (!tmp[no])
        tmp[no] = hash_open(no, "w");
    fprintf(tmp[no], "%s\n", url);
}

/**
 * 获取指定哈希小文件
 * @param 哈希值
 * @return 文件指针
 */
FILE *hash_get(int no) {
    if (tmp[no]) {
        fclose(tmp[no]);
        tmp[no] = hash_open(no, "r");
    }
    return tmp[no];
}

/**
 * 从文件中读入一个 url
 * @param 文件指针
 * @return url内容
 */
const char *fgeturl(FILE *f) {
    int ch = fgetc(f), len = 0;
    static char url[MAX_URL_LEN];
    do {
        if (url_char(ch)) break;
        if (ch == EOF) return NULL;
        ch = fgetc(f);
    } while (true);
    do {
        url[len++] = ch;
        if (len >= MAX_URL_LEN) {
            fprintf(stderr, "too long url");
            exit(1);
        }
    } while (url_char(ch = fgetc(f)));
    url[len] = 0;
    return url;
}

/**
 * 删除临时创建的哈希小文件
 */
void remove_tmp() {
    for (int i = 0; i < HASH_MOD; i++)
        if (tmp[i]) {
            char buf[20];
            sprintf(buf, "HASH_%d", i);
            remove(buf);
        }
}

/**
 * 计算字符串哈希
 * @param str 字符串
 * @return 哈希值
 */
int string_hash(const char *str) {
    int res = 0, len = strlen(str);
    for (int i = 0; i < len; i++) {
        res = res * 31 + str[i];
        res %= HASH_MOD;
    }
    return res;
}
