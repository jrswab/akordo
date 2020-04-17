# The preferred contribution approach:
1. Clone
2. Make a branch named with your fix or new feature
3. Code
4. Run the unit tests
5. Update tests as needed
6. Switch back to master
7. `git pull origin master`
8. `git merge [your branch name]` 
9. Resolve any conflicts
10. Use `git send-email --to="~jrswab/akordo-dev@lists.sr.ht" HEAD^` to create and send a patch
 
Please refer to the sourcehut [Email Etiquettes](https://man.sr.ht/lists.sr.ht/etiquette.md) when sending patches.
Using `git send-email` will take care of most of the etiquette for you.

To learn how to send patches with `git send-email` check out the sourcehut tutorial at [git-send-email.io](https://git-send-email.io/)