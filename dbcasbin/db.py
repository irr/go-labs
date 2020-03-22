import os

USERS = 100000

RULES = [
    "/billing/gym/{}/report/{}",
    "/payout/gym/{}/report/{}",
]

ACTIONS = [
    ["(billing.view_report)|(billing.edit_report)", "(billing.view_report)"],
    ["(payout.view_report)|(payout.edit_report)", "(payout.view_report)"],
]

POLICIES = [
    "p, nobody, /billing/gym/:gid/report/:rid, validate",
    "p, nobody, /payout/gym/:gid/report/:rid, validate",
    "p, nobody, /profile, validate",
    "p, role:users, /profile, (profile.view)|(profile.edit)",
    "g, role:managers, role:users",
]

BILLING = ["(billing.view_report)|(billing.edit_report)", "(billing.view_report)"]
PAYOUT = ["(payout.view_report)|(payout.edit_report)", "(payout.view_report)"]


for p in POLICIES:
    print(p)


for u in range(USERS):
    if u < (USERS / 4):
        print(f"g, user{u+1}, role:managers")
    else:
        print(f"g, user{u+1}, role:users")
    for i, r in enumerate(RULES):
        a, b = f"g-{u+1}", f"r-{u+1}"
        s = f"p, user{u+1}, {r.format(a, b)}, {ACTIONS[i][u % 2]}"
        print(s)
