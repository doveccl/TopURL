#include <queue>
#include <vector>
#include <unordered_map>
#include "util.h"

#define TOP_K 100

using namespace std;

// url 文件
FILE *f_url;

// 小文件 url 计数
unordered_map<string, int> M;

typedef pair<string, int> PSI;
auto cmp = [](PSI a, PSI b) { return a.second > b.second; };
// 保存最终结果的优先队列
priority_queue<PSI, vector<PSI>, decltype(cmp)> Q(cmp);

int main(int argc, char **argv) {
    if (argc < 2) {
        // 检查参数个数
        fprintf(stderr, "url file is required");
        exit(1);
    }
    if (!(f_url = fopen(argv[1], "r"))) {
        // 打开 url 文件
        fprintf(stderr, "fail to open url file");
        exit(1);
    }
    while (auto s = fgeturl(f_url)) {
        // 将 url 按照哈希结果写入小文件
        hash_add(string_hash(s), s);
    }
    int qlen = 0; // 答案队列长度
    for (int i = 0; i < HASH_MOD; i++) {
        // 遍历哈希小文件
        FILE *f_hash = hash_get(i);
        if (!f_hash) continue;
        M.clear();
        while (auto s = fgeturl(f_hash)) {
            // url 计数统计
            M[s] = M[s] + 1;
        }
        for (auto p: M) {
            // 依次加入优先队列，控制队列长度
            Q.push(p), qlen++;
            while (qlen > TOP_K)
                Q.pop(), qlen--;
        }
    }
    // 输出 TOP_K 的 url 和对应的计数
    printf("CNT\tURL\n");
    while (!Q.empty()) {
        auto p = Q.top(); Q.pop();
        printf("%d\t%s\n", p.second, p.first.c_str());
    }
    // 删除小文件
    remove_tmp();
    return 0;
}
