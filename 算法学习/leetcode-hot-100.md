[TOC]

# LeeCode HOT 100

> https://leetcode.cn/problem-list/2cktkvj/

### 2.两数相加

由于输入的两个链表都是逆序存储数字的位数的，因此两个链表中同一位置的数字可以直接相加。

我们同时遍历两个链表，逐位计算它们的和，并与当前位置的进位值相加。具体而言，如果当前两个链表处相应位置的数字为 n1,n2，进位值为 carry，则它们的和为 `n1+n2+ carry`；其中，答案链表处相应位置的数字为 `(n1+n2+ carry) %10`，而新的进位值为 `(n1+n2+carry) / 10`

如果两个链表的长度不同，则可以认为长度短的链表的后面有若干个 0 。

此外，如果链表遍历结束后，有 carry>0，还需要在答案链表的后面附加一个节点，节点的值为 carry。

### [3.无重复字符的最长子串](https://leetcode.cn/problems/longest-substring-without-repeating-characters/)

滑动窗口算法

### [11. 盛最多水的容器](https://leetcode.cn/problems/container-with-most-water/)

这里需要注意的是双指针要怎么去移动。很显然短的那条边移动

### [17. 电话号码的字母组合](https://leetcode.cn/problems/letter-combinations-of-a-phone-number/)

回溯法。参考题解：[电话号码的字母组合](https://leetcode.cn/problems/letter-combinations-of-a-phone-number/solution/leetcode-17-letter-combinations-of-a-phone-number-/)
